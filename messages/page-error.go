package messages

type Message struct {
	ErrorType    string
	ErrorHeader  string
	ErrorMessage string
	Redirect     string
}

var Error404 = Message{
	ErrorType:    "404",
	ErrorHeader:  "Oops! Page not found",
	ErrorMessage: "The page you're looking for doesn't exist or has been moved.",
	Redirect:     "/",
}

var Error403 = Message{
	ErrorType:    "403",
	ErrorHeader:  "Access denied",
	ErrorMessage: "Your are not allowed access to access this page ",
}

var Error400 = Message{
	ErrorType:    "400",
	ErrorHeader:  "Bad Request",
	ErrorMessage: "Bad Request page cannot be retrieved",
	Redirect:     "/",
}
var Error500 = Message{
	ErrorType:    "500",
	ErrorHeader:  "Internal Server Error",
	ErrorMessage: "Internal Server Error page cannot be retrieved",
	Redirect:     "/",
}
