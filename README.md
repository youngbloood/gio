[![codecov](https://codecov.io/gh/youngbloood/gio/branch/master/graph/badge.svg)](https://codecov.io/gh/youngbloood/gio)
# gio
goroutine multiplex .


# Usage

```

func main(){
    result := make(Result, 5)
    m := gio.NewMaster(5, result)
    m.Run()
    defer m.Stop()
    var jobs = newJobs(100)
    m.Push(jobs...)
    runtime.Gosched()
    log.Println("result = ", result)
}


func newJobs(length int) (jobs []gio.Jober) {
    for i := 0; i < length; i++ {
        jobs = append(jobs, job{i})
    }
    return
}

// Next implement by yourself.

// Job is custom job.
type job struct {
    A int
}

// implement the gio.Jober interface
func (j job) Run(result gio.Saver) error {
    // ...

    // save the result
    return result.Save(fmt.Sprintf("job%d", j.A), j.A)
}

var rwmux sync.RWMutex

// Result is custom result
type Result map[string]interface{}

// implement the gio.Saver interface.(you can handle it with sync or async)
func (r Result) Save(v ...interface{}) error {
    rwmux.Lock()
    r[v[0].(string)] = v[1]
    rwmux.Unlock()
    return nil
}

```


