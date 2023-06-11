package common

func setSuccess(err error, b *bool) error {
	*b = err == nil
	return err
}

// ReadingService adapts ReadingList for RPC
type ReadingService struct {
	ReadingList
}

func (r *ReadingService) AddCourse(b Course, success *bool) error {
	return setSuccess(r.ReadingList.AddCourse(b), success)
}

func (r *ReadingService) RemoveCourse(isbn string, success *bool) error {
	return setSuccess(r.ReadingList.RemoveCourse(isbn), success)
}

func (r *ReadingService) GetProgress(isbn string, pages *int) (err error) {
	*pages, err = r.ReadingList.GetProgress(isbn)
	return err
}

type Progress struct {
	ISBN  string
	Pages int
}

func (r *ReadingService) SetProgress(p Progress, success *bool) error {
	return setSuccess(r.ReadingList.SetProgress(p.ISBN, p.Pages), success)
}

func (r *ReadingService) AdvanceProgress(p Progress, success *bool) error {
	return setSuccess(r.ReadingList.AdvanceProgress(p.ISBN, p.Pages), success)
}
