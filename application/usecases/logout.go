package usecases

func (uc *AuthUseCase) Logout(userID, token string) error {
	if uc.redisClient == nil {
		return nil
	}
	return uc.redisClient.DeleteSession(token)
}
