.PHONY: git
git:
	go build -o main && git add . && git commit -m "Update content" && git push origin master