Version, Health and Metrics helpers
===================

How to Use it
-------------

```
rest.RegistryVersion(&struct {
		rest.Version
		Custom string
	}{
		Version: rest.Version{
			GeneralInfo: rest.GeneralInfo{
				ServiceName:     "ServiceName",
				ServiceProvider: "ServiceProvider",
				ServiceVersion:  "1.0.0",
				Name:            "SolutionName",
			},
			Type:                 "VersionType",
			SupportedAPIVersions: []string{"1"},
			BuildNumber:          "BuildNumber",
			BuildCommitSHA:       "BuildCommitSHA",
			Repository:           "Repository",
		},
		Custom: "Custom field",
	})
```
and
```
rest.RegistryHealth(&struct {
		rest.Health
		Custom string
	}{
		Health: rest.Health{
			GeneralInfo: rest.GeneralInfo{
				ServiceName:     "ServiceName",
				ServiceProvider: "ServiceProvider",
				ServiceVersion:  "1.0.0",
				Name:            "SolutionName",
			},
			ConnMethods: []rest.Statuser{},
			ListenURL: ":12124",
		},
		Custom: "Custom field",
	})
```
and 

```
rest.RegistryMetrics(&struct {
		rest.Metrics
	}{
	Metrics: rest.Metrics{
		GeneralInfo: rest.GeneralInfo{
			ServiceName:     "ServiceName",
			ServiceProvider: "ServiceProvider",
			ServiceVersion:  "1.0.0",
			Name:            "SolutionName",
		},
	}}, []rest.MetricsConfig{
	{
        Name:     alertLatency,
        DataType: rest.DataTypeDecimal,
        Range:    rest.RangeTypeInfinity,
        Unit:     rest.UnitTypeNumbers,
    }},
)

// increment metrics by value
rest.Add("metricName", 1)

or

// set metrics value
rest.Set("metricName", 10)

or 

type Metrics struct {
	rest.Metrics
}
func (m Metrics) Update() error {
    rest.Set("metricName", 10)
    return nil
}
```

then create routes for your service

```
http.HandleFunc("/version", rest.HandlerVersion)
http.HandleFunc("/health", rest.HandlerHealth)
http.HandleFunc("/metrics", rest.HandlerMetrics)
```
