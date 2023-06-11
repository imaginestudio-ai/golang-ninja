package persistence

type DatabaseHandler interface {
	AddUser(User) ([]byte, error)
	AddEvent(Event) ([]byte, error)
	AddCourseingForUser([]byte, Courseing) error
	AddLocation(Location) (Location, error)
	FindUser(string, string) (User, error)
	FindCourseingsForUser([]byte) ([]Courseing, error)
	FindEvent([]byte) (Event, error)
	FindEventByName(string) (Event, error)
	FindAllAvailableEvents() ([]Event, error)
	FindLocation(string) (Location, error)
	FindAllLocations() ([]Location, error)
}
