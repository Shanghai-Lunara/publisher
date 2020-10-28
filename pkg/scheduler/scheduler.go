package scheduler

type Scheduler struct {
	Namespaces []Namespace
}

type Namespace struct {
	Groups []Group
}

type Group struct {

}
