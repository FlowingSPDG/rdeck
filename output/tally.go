package output

type Tally interface {
	Active() error
	Preview() error
	Inactive() error
	// TODO: 初期化/終了処理...?
}
