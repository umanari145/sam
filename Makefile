.PHONY: build


debug:
	lambdaSample/build.sh && sam local invoke LambdaSampleFunction --event event_api.json	


build:
	sam build
