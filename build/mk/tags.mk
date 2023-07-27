recreate-tag: delete-tag create-tag

create-tag:
	@echo "Tagging version $(VERSION)"
	@git tag -a $(VERSION) -m "Release $(VERSION)"
	@git push origin $(VERSION)

delete-tag:
	@echo "Deleting tag $(VERSION)"
	@git tag -d $(VERSION)
	@git push origin --delete $(VERSION)