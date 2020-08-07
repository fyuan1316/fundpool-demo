package v1alpha1

func (s SupplyStatus) IsReachFinal() bool {
	if s.Phase != "" {
		return true
	}
	return false
}

func (s *SupplyStatus) MarkFailed() {
	s.setState(Failed)
}
func (s *SupplyStatus) MarkSuccessed() {
	s.setState(Succeeded)
}
func (s *SupplyStatus) setState(state PhaseState) {
	s.Phase = state
}
