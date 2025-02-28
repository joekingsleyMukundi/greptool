package worklist

type Entry struct {
	Path string
}

type Worklist struct {
	job chan Entry
}

func (w *Worklist) Add(work Entry) {
	w.job <- work
}
func (w *Worklist) Next() Entry {
	j := <-w.job
	return j
}

func New(buffSize int) Worklist {
	return Worklist{make(chan Entry, buffSize)}
}
func NewJob(path string) Entry {
	return Entry{path}
}

func (w *Worklist) Finalize(numWorkers int) {
	for i := 0; i < numWorkers; i++ {
		w.Add(Entry{})
	}
}
