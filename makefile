.PHONY: git
#git:
#	git add . && git commit -m "Update content" && git push

git:
	git config user.name "giftia" && \
    git config user.email "hch20041214sr@qq.com" && \
    git add . && \
    git commit -m "Update content" && \
    git push

# git push origin master