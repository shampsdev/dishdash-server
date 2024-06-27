package entities

type Vote struct {
	options   []int
	voteCount int
	onEnd     func([]int)
	onVote    func(*Vote, []int)
	ended     bool
}

func NewVote(numOptions int, onVote func(*Vote, []int), onEnd func([]int)) *Vote {
	return &Vote{
		options: make([]int, numOptions),
		onEnd:   onEnd,
		onVote:  onVote,
		ended:   false,
	}
}

func (v *Vote) Vote(option int) {
	if option >= 0 && option < len(v.options) {
		v.options[option]++
		v.voteCount++
		v.onVote(v, v.options)
	}
}

func (v *Vote) FinalizeVote() {
	if v.onEnd != nil && !v.ended {
		v.onEnd(v.options)
		v.ended = true
	}
}
