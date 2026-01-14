package provider

// type ProvidersStore interface {
// 	Create(ctx context.Context, provider *Provider) error
// 	Update(ctx context.Context, provider *Provider) error
// 	Delete(ctx context.Context, id uuid.UUID) error
// 	GetDetailsByID(ctx context.Context, id uuid.UUID) (*Provider, error)
// 	GetAllAvailable(ctx context.Context) ([]*Provider, error)
// 	List(ctx context.Context) ([]*Provider, error)
// }

// type PostgreSQLProvidersStore struct {
// 	q *generated.Queries
// }

// func NewProvidersStore(provider db.DBProvider) ProvidersStore {
// 	return &PostgreSQLProvidersStore{
// 		q: generated.New(provider.GetConn()),
// 	}
// }

// func (s *PostgreSQLProvidersStore) Create(ctx context.Context, provider *Provider) error {
// 	params := generated.CreateProviderParams{
// 		ID:        provider.ID,
// 		Name:      provider.Name,
// 		CreatedAt: dbutils.ConvertTimeToPgtypeTimestamp(time.Now()),
// 	}

// 	_, err := s.q.CreateProvider(
// 		ctx,
// 		params,
// 	)

// 	return err
// }

// func (s *PostgreSQLProvidersStore) Update(ctx context.Context, provider *Provider) error {
// 	params := generated.UpdateProviderParams{
// 		ID:   provider.ID,
// 		Name: provider.Name,
// 		ModifiedAt: dbutils.ConvertNullableTimeToPgtypeTimestamp(
// 			provider.ModifiedAt,
// 		),
// 	}

// 	err := s.q.UpdateProvider(
// 		ctx,
// 		params,
// 	)

// 	return err
// }

// func (s *PostgreSQLProvidersStore) Delete(ctx context.Context, id uuid.UUID) error {
// 	params := generated.DeleteProviderParams{
// 		ID:         id,
// 		ModifiedAt: dbutils.ConvertTimeToPgtypeTimestamp(time.Now()),
// 		DeletedAt:  dbutils.ConvertTimeToPgtypeTimestamp(time.Now()),
// 	}

// 	err := s.q.DeleteProvider(ctx, params)

// 	return err
// }

// func (s *PostgreSQLProvidersStore) GetDetailsByID(ctx context.Context, id uuid.UUID) (*Provider, error) {
// 	row, err := s.q.GetProviderDetailsByProviderID(ctx, id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	modifiedAt, err := dbutils.ConvertPgtypeTimestampToTimePointer(row.ModelProvider.ModifiedAt)
// 	if err != nil {
// 		return nil, err
// 	}

// 	deletedAt, err := dbutils.ConvertPgtypeTimestampToTimePointer(row.ModelProvider.DeletedAt)
// 	if err != nil {
// 		return nil, err
// 	}

// 	provider := FromDatabase(
// 		row.ModelProvider.ID,
// 		row.ModelProvider.Name,
// 		row.ModelProvider.CreatedAt.Time,
// 		modifiedAt,
// 		deletedAt,
// 	)

// 	return provider, nil
// }

// func (s *PostgreSQLProvidersStore) GetAllAvailable(ctx context.Context) ([]*Provider, error) {
// 	rows, err := s.q.GetAllProviders(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var providers []*Provider

// 	for _, row := range rows {
// 		modifiedAt, err := dbutils.ConvertPgtypeTimestampToTimePointer(row.ModelProvider.ModifiedAt)
// 		if err != nil {
// 			return nil, err
// 		}

// 		deletedAt, err := dbutils.ConvertPgtypeTimestampToTimePointer(row.ModelProvider.DeletedAt)
// 		if err != nil {
// 			return nil, err
// 		}

// 		provider := FromDatabase(
// 			row.ModelProvider.ID,
// 			row.ModelProvider.Name,
// 			row.ModelProvider.CreatedAt.Time,
// 			modifiedAt,
// 			deletedAt,
// 		)

// 		providers = append(providers, provider)
// 	}

// 	return providers, nil
// }

// func (s *PostgreSQLProvidersStore) List(ctx context.Context) ([]*Provider, error) {
// 	rows, err := s.q.GetAllProviders(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var providers []*Provider

// 	for _, row := range rows {
// 		modifiedAt, err := dbutils.ConvertPgtypeTimestampToTimePointer(row.ModelProvider.ModifiedAt)
// 		if err != nil {
// 			return nil, err
// 		}

// 		deletedAt, err := dbutils.ConvertPgtypeTimestampToTimePointer(row.ModelProvider.DeletedAt)
// 		if err != nil {
// 			return nil, err
// 		}

// 		provider := FromDatabase(
// 			row.ModelProvider.ID,
// 			row.ModelProvider.Name,
// 			row.ModelProvider.CreatedAt.Time,
// 			modifiedAt,
// 			deletedAt,
// 		)

// 		providers = append(providers, provider)
// 	}

// 	return providers, nil
// }
