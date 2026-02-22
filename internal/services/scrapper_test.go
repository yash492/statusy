package applications

func (t *TestSuite) TestScrapper() {
	StartScrapper(Deps{
		Logger:  t.Logger,
		ReadDB:  t.TestDb,
		WriteDB: t.TestDb,
	})
}
