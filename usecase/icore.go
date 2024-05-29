package usecase

import (
	core_actor "filmoteka/usecase/actors"
	core_films "filmoteka/usecase/films"
	core_profiles "filmoteka/usecase/profiles"
	core_sessions "filmoteka/usecase/sessions"
)

type ICore interface {
	core_films.IFilms
	core_actor.IActors
	core_profiles.IProfiles
	core_sessions.ISessions
}
