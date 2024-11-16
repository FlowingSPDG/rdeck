package output

type Out interface {
	Name() string
}

type Digital interface {
	Out

	On() error
	Off() error
	// TODO: 初期化/終了処理...?
}
