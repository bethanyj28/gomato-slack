build:
	docker build -t bethanyj28/gomato-slack .
run:
	docker run --rm -p 8080:8080 bethanyj28/gomato-slack
env:
	touch environment.env
	echo "SLACK_VERIFICATION_TOKEN=<Your Slack verification token>" >> environment.env
