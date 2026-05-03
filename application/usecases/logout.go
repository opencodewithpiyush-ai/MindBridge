package usecases

func (uc *AuthUseCase) Logout(userID, jti string) error {
	if uc.redisClient == nil {
		return nil
	}
	return uc.redisClient.DeleteSession(jti)
}
