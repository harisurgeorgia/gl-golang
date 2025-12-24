package messages

type Message struct {
	ErrorType    string
	ErrorHeader  string
	ErrorMessage string
}

var Error404 = Message{
	ErrorType:    "404",
	ErrorHeader:  "Oops! Page not found",
	ErrorMessage: "The page you're looking for doesn't exist or has been moved.",
}

var Error403 = Message{
	ErrorType:    "403",
	ErrorHeader:  "Access denied",
	ErrorMessage: "Your are not access to access this page ",
}
