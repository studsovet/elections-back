package controllers

/*
type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	t, err := db.LoginCheck(input.Username, input.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username or password is incorrect."})
		return
	}

	token.SetTokenCookie(c, t)
	c.JSON(http.StatusOK, gin.H{"token": t})
}

type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	u := db.User{}
	u.Username = input.Username
	u.Password = input.Password
	u.IsObserver = false
	_, err := u.SaveUser()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "registration success"})
}

func GetCurrentUser(c *gin.Context) (db.User, error) {
	user_id, err := token.ExtractTokenID(c)
	if err != nil {
		return db.User{}, err
	}

	u, err := db.GetUserByID(user_id)

	if err != nil {
		return db.User{}, err
	}
	return u, nil
}
*/
