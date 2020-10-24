package jobs

// Фоновый процесс, исполняемый в горутине
type Job interface {
	Run()
}
