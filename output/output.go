package output

type Analog interface {
	On() error
	Off() error
	// TODO: 初期化/終了処理...?
}
