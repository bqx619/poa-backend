package biz

import (
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"poa-service/app/analysis/internal/conf"
	"poa-service/app/analysis/internal/data/ent"
	"poa-service/app/analysis/internal/data/ent/zeeklog"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/panjf2000/ants/v2"
	"github.com/pkg/errors"
)

var (
	dirNameReplace  = strings.NewReplacer("{Y}", "2006", "{M}", "01", "{D}", "02")
	fileNameReplace = strings.NewReplacer("type", "%s", "begin_hour", "%s", "end_hour", "%s")
	gp, _ = ants.NewPool(5000)
)

type Opinion struct {
	Ts        int
	Uid       string
	OrigH     string
	OrigP     int32
	RespH     string
	RespP     int32
	Domain    string
	Host      string
	Proto     string
	Uri       string
	UserAgent string
	Title     string
	Content   string
	Body      io.Reader
}

type OpinionRepo interface {
	// mysql
	ListPendingOpinion(ctx context.Context) ([]*ent.Opinion, error)
	CreateOpinion(ctx context.Context, opinion *Opinion) error
	UpdateAnalysisStatus(ctx context.Context, id int, status interface{}) error
}

type OpinionUseCase struct {
	repo OpinionRepo
	log  *log.Helper
}

type ZeekLogRepo interface {
	GetZeekLog(ctx context.Context, md5 string) (*ent.ZeekLog, error)
	GetZeekLogByFile(ctx context.Context, dir, fileName string) (*ent.ZeekLog, error)
	GetZeekLogById(ctx context.Context, id int) (*ent.ZeekLog, error)
	GetLastZeekLog(ctx context.Context) (*ent.ZeekLog, error)
	CreateZeekLog(ctx context.Context, zeekLog *ZeekLog) error
	SetRunningCount(ctx context.Context, count int) error
	GetRunningCount(ctx context.Context) (int, error)
	SetRunningTag(ctx context.Context, tag int) error
	GetRunningTag(ctx context.Context) (int, error)
	Count(ctx context.Context) (int, error)
	ListZeekLogs(ctx context.Context, sort string, sortType string, pageSize int, current int, id int, dir string, day int64) ([]*ent.ZeekLog, int, error)
}

type ZeekLog struct {
	Dir       string
	FileName  string
	Type      zeeklog.Type
	BeginHour int8
	EndHour   int8
	Md5       string
	LogDay    int
}

type RefactorUseCase struct {
	opinionRepo OpinionRepo
	zeekLogRepo ZeekLogRepo
	conf        *conf.ZeekLog
	log         *log.Helper
}

func NewRefactorUseCase(c *conf.ZeekLog, opinionRepo OpinionRepo, zeekLogRepo ZeekLogRepo, logger log.Logger) *RefactorUseCase {
	return &RefactorUseCase{
		opinionRepo: opinionRepo,
		zeekLogRepo: zeekLogRepo,
		log:         log.NewHelper(logger),
		conf:        c,
	}
}

func (ref *RefactorUseCase) Status(ctx context.Context) ([]int, error) {
	fileInfos, err := ioutil.ReadDir(ref.conf.Dir)
	if err != nil {
		return nil, errors.WithMessage(err, "open zeek log dir failed")
	}
	allDataCount := 0
	dirFormat := dirNameReplace.Replace(ref.conf.DirNameFormat)
	for _, fileInfo := range fileInfos {
		if !fileInfo.IsDir() {
			continue
		}
		_, err := time.Parse(dirFormat, fileInfo.Name())
		if err != nil {
			continue
		}
		allDataCount++
	}
	refactoredCount, err := ref.zeekLogRepo.Count(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "count refactored failed")
	}
	runningCount, err := ref.zeekLogRepo.GetRunningCount(ctx)
	if err == redis.Nil || err != nil {
		runningCount = 0
	}
	waitingCount := allDataCount - refactoredCount - runningCount
	return []int{waitingCount, runningCount, 0, refactoredCount}, nil
}

func (ref *RefactorUseCase) List(ctx context.Context, sort string, sortType string, pageSize int, current int, id int, dir string, day int64) ([]*ent.ZeekLog, int, error) {
	return ref.zeekLogRepo.ListZeekLogs(ctx, sort, sortType, pageSize, current, id, dir, day)
}

func (ref *RefactorUseCase) Get(ctx context.Context, id int) (*ent.ZeekLog, error) {
	return ref.zeekLogRepo.GetZeekLogById(ctx, id)
}

func (ref *RefactorUseCase) Run(ctx context.Context) error {
	ref.log.Info("begin refactor running ...")
	tag, err := ref.zeekLogRepo.GetRunningTag(ctx)
	if err != nil {
		ref.log.Error(err)
		tag = 0
	}
	if tag == 1 {
		ref.log.Info("tag = 1")
		return nil
	}
	fileBlockPool := sync.Pool{New: func() interface{} { return make([]byte, 1000*1024) }}
	opinionPool := sync.Pool{New: func() interface{} { return new(Opinion) }}
	zeekLogPool := sync.Pool{New: func() interface{} { return new(ZeekLog) }}

	_ = ref.zeekLogRepo.SetRunningTag(ctx, 1)

	defer func(zeekLogRepo ZeekLogRepo, ctx context.Context, tag int) {
		_ = zeekLogRepo.SetRunningTag(ctx, tag)
	}(ref.zeekLogRepo, ctx, 0)

	for {
		// 工作目录列表
		fileInfos, err := ioutil.ReadDir(ref.conf.Dir)
		if err != nil {
			return errors.WithMessage(err, "open zeek log dir failed")
		}

		dirFormat := dirNameReplace.Replace(ref.conf.DirNameFormat)

		type dir struct {
			name string
			day  time.Time
		}

		dirs := make([]dir, 0)

		for _, fileInfo := range fileInfos {
			if !fileInfo.IsDir() {
				continue
			}
			parse, err := time.Parse(dirFormat, fileInfo.Name())
			if err != nil {
				ref.log.Warnf("parse dir err=%+v", err)
				continue
			}
			dirs = append(dirs, dir{name: fileInfo.Name(), day: parse})
		}

		sort.Slice(dirs, func(i, j int) bool {
			return dirs[i].day.Before(dirs[i].day)
		})

		// 上次结束的地方开始
		lastZeekLog, err1 := ref.zeekLogRepo.GetLastZeekLog(ctx)

		for _, d := range dirs {
			if err1 == nil && d.day.Before(time.UnixMicro(int64(lastZeekLog.LogDay))) {
				continue
			}
			_, err := ref.zeekLogRepo.GetZeekLogByFile(context.TODO(), ref.conf.Dir, d.name)
			if !ent.IsNotFound(err) {
				ref.log.Infof("zeek log had record, err=%+v", err)
				continue
			}
			dirFiles, err2 := ioutil.ReadDir(filepath.Join(ref.conf.Dir, d.name))
			if err2 != nil {
				ref.log.Errorf("io err=%+v", err2)
				continue
			}
			for i := range dirFiles {
				if strings.HasPrefix(dirFiles[i].Name(), "http.") && strings.HasSuffix(dirFiles[i].Name(), ".log.gz") {

					hourSplit := strings.Split(strings.Split(dirFiles[i].Name(), ".")[1], "-")
					beginHour, _ := time.Parse("15:04:05", hourSplit[0])
					endHour, _ := time.Parse("15:04:05", hourSplit[1])
					zeekLog := zeekLogPool.Get().(*ZeekLog)
					zeekLog.Dir = filepath.Join(ref.conf.Dir, d.name)
					zeekLog.FileName = dirFiles[i].Name()
					zeekLog.Type = zeeklog.TypeHTTP
					zeekLog.BeginHour = int8(beginHour.Hour())
					zeekLog.EndHour = int8(endHour.Hour())
					zeekLog.LogDay = int(d.day.UnixMicro())
					err = ref.zeekLogRepo.CreateZeekLog(context.TODO(), zeekLog)
					if err != nil {
						ref.log.Errorf("save zl=%+v err=%+v", zeekLog, err)
						continue
					}
					zeekLogPool.Put(zeekLog)

					ref.readFile(
						ctx,
						filepath.Join(ref.conf.Dir, d.name, dirFiles[i].Name()),
						&fileBlockPool, &opinionPool)
				}
			}

		}
		time.Sleep(ref.conf.Interval.AsDuration())
	}
}

func (ref *RefactorUseCase) ListUnAnalysis(ctx context.Context) ([]*ent.Opinion, error) {
	return ref.opinionRepo.ListPendingOpinion(ctx)
}

func (ref *RefactorUseCase) UpdateAnalysisStatus(ctx context.Context, id int, status string) error {
	return ref.opinionRepo.UpdateAnalysisStatus(ctx, id, status)
}

func (ref *RefactorUseCase) readFile(ctx context.Context, fileName string, fileBlockPool *sync.Pool, opinionPool *sync.Pool) {
	var wg sync.WaitGroup

	gp.Reboot()

	file, err := os.Open(fileName)
	if err != nil {
		ref.log.Warnf("open file=%s error=%+v", fileName, err)
		return
	}

	r, err1 := gzip.NewReader(file)

	defer func() {
		_ = r.Close()
		_ = file.Close()
		gp.Release()
	}()

	if err1 == io.EOF {
		ref.log.Warnf("parse gzip file=%s EOF", fileName)
		return
	}

	if err1 != nil {
		ref.log.Warnf("parse gzip file=%s error=%+v", fileName, err)
		return
	}

	lastSurplus := ""
	//allLines := make([]*[30]string, 0)
	dc := make(chan *Opinion, 200)
	//ch := make(chan struct{}, 10000)
	_ = gp.Submit(func() {
		for {
			select {
			case op := <-dc:
				fmt.Println("get ...")
				err = ref.opinionRepo.CreateOpinion(context.TODO(), op)
				if err != nil {
					ref.log.Errorf("create opinion op=%+v error=%+v", op, err)
				}
				opinionPool.Put(op)
			}
		}
	})
	detailLines := 0
	for {
		buf := fileBlockPool.Get().([]byte)
		read, err := r.Read(buf)
		buf = buf[:read]
		if read == 0 {
			if err == io.EOF {
				break
			}
			if err != nil {
				ref.log.Warnf("%+v", err)
				break
			}
			return
		}
		lines := lastSurplus + string(buf)

		if lastSurplus != "" {
			lastSurplus = ""
		}

		for _, line := range strings.Split(lines, "\n") {
			if strings.HasPrefix(line, "#") {
				continue
			}
			splitLine := strings.Split(line, "\t")
			if len(splitLine) != 30 {
				lastSurplus = line
			} else {

				//ch <- struct{}{}
				wg.Add(1)
				detailLines ++
				_ = gp.Submit(func() {
					err := ref.check(opinionPool, splitLine, dc)
					if err != nil {
						//fmt.Printf("err=%v \n", err)
					}
					wg.Done()
				})
				//go ref.check(opinionPool, splitLine, &wg, dc, ch)

				//checked, err := ref.check(opinionPool, splitLine)
				//if err == nil {
				//	ch <- struct{}{}
				//	fmt.Println("get line ...", k)
				//	wg.Add(1)
				//	runningCount++
				//	//go ref.readLine(checked, &wg, dc)
				//	go func() {
				//		defer wg.Done()
				//		fmt.Println("get line ...", k)
				//		err = ref.opinionRepo.CreateOpinion(context.TODO(), checked)
				//		if err != nil {
				//			ref.log.Errorf("create opinion op=%+v error=%+v", checked, err)
				//		}
				//		opinionPool.Put(checked)
				//		//dc <- checked
				//	}()
				//}
			}
		}
	}

	wg.Wait()
	_ = ref.zeekLogRepo.SetRunningCount(ctx, gp.Running())
	fmt.Println("Finish processing a file, detail lines:", detailLines)
}

func GetStrCn(str string) (cnStr string) {
	r := []rune(str)
	var strSlice []string
	for i := 0; i < len(r); i++ {
		if r[i] <= 40869 && r[i] >= 19968 {
			cnStr = cnStr + string(r[i])
			strSlice = append(strSlice, cnStr)
		}
	}
	return
}

func (ref *RefactorUseCase) check(opinionPool *sync.Pool, line []string,dc chan <- *Opinion) error {
	if line[7] != "GET" {
		return fmt.Errorf("not GET")
	}
	timestamp, err := strconv.Atoi(strings.ReplaceAll(line[0], ".", ""))
	if err != nil {
		return fmt.Errorf("timestamp format error=%+v", err)
	}

	origP, err := strconv.Atoi(line[3])
	if err != nil {
		return fmt.Errorf("origP format error=%+v", err)
	}
	respP, err := strconv.Atoi(line[5])
	if err != nil {
		return fmt.Errorf("origP format error=%+v", err)
	}
	requestURL := fmt.Sprintf("http://%s:%d%s", line[8], respP, line[9])
	parseURL, err := url.Parse(requestURL)
	if err != nil {
		return fmt.Errorf("parse error=%+v", err)
	}
	if line[8] == "-" || !strings.Contains(line[8], ".") {
		return fmt.Errorf("parse domain url err")
	}
	response, err := http.Get(requestURL)

	opinion := opinionPool.Get().(*Opinion)
	opinion.Ts = timestamp
	opinion.Uid = line[1]
	opinion.OrigH = line[2]
	opinion.OrigP = int32(origP)
	opinion.RespH = line[4]
	opinion.RespP = int32(respP)
	opinion.Proto = "http"
	opinion.UserAgent = line[12]
	opinion.Uri = line[9]
	opinion.Host = line[8]
	splitURL := strings.Split(parseURL.Hostname(), ".")
	opinion.Domain = splitURL[len(splitURL)-2] + "." + splitURL[len(splitURL)-1]
	if err != nil {
		return fmt.Errorf("request err=%+v", err)
	}
	defer response.Body.Close()

	if response.StatusCode == 200 {
		opinion.Body = response.Body
		doc, err := goquery.NewDocumentFromReader(opinion.Body)
		if err != nil {
			return fmt.Errorf("query html err=%+v", err)
		}

		cnBody := GetStrCn(doc.Find("body").Text())
		if len(cnBody) < 10 {
			return fmt.Errorf("short html")
		}

		opinion.Title = doc.Find("title").Text()
		opinion.Content = cnBody
	}

	fmt.Println("send data ...")
	dc <- opinion
	//<- ch
	return nil
}

func (ref *RefactorUseCase) readLine(op *Opinion, wg *sync.WaitGroup, dc chan<- *Opinion) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recover ...:", r)
		}
		wg.Done()
	}()

	doc, err := goquery.NewDocumentFromReader(op.Body)
	if err != nil {
		return
	}

	cnBody := GetStrCn(doc.Find("body").Text())
	if len(cnBody) < 50 {
		return
	}

	op.Title = doc.Find("title").Text()
	op.Content = cnBody

	fmt.Println("send ...")
	dc <- op

	return
}
