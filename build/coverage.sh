cat test_cov/cover.out.tmp | grep -v ".pb.go" >cover.out
go tool cover -func cover.out
