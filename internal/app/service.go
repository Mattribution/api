package app

type Service struct {
	tracksDAO TracksDAO
	kpisDAO   KpisDAO
}

func (s Service) NewTrack() error {

}
