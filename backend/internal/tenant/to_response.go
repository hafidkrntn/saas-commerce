package tenant

// ToResponse converts a Tenant model to its API response representation.
func (t *Tenant) ToResponse() TenantResponse {
	return TenantResponse{
		ID:                 t.ID,
		Name:               t.Name,
		Slug:               t.Slug,
		Logo:               t.Logo,
		PlanID:             t.PlanID,
		SubscriptionStatus: t.SubscriptionStatus,
		TrialEndsAt:        t.TrialEndsAt,
		CreatedAt:          t.CreatedAt,
	}
}

// ToResponse converts a Plan model to its API response representation.
func (p *Plan) ToResponse() PlanResponse {
	return PlanResponse{
		ID:          p.ID,
		Name:        p.Name,
		Code:        p.Code,
		Price:       p.Price,
		Description: p.Description,
	}
}
