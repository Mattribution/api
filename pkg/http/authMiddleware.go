package http

// Define our struct
type authMiddleware struct {
	tokenUsers map[string]int
	h          Handler
}

func newAuthMiddleware(h Handler) authMiddleware {
	amw := authMiddleware{
		h: h,
	}
	amw.Populate()
	return amw
}

// Initialize it somewhere
func (amw *authMiddleware) Populate() {
	amw.tokenUsers["00000000"] = 1
}
