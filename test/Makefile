DIRS=ws

tests:
	cd usersim;make

test:
	@echo "FUNCTIONAL TEST Report" > testreport.txt
	@echo "Report Start: \c" >> testreport.txt
	@date >> testreport.txt
	@echo >> testreport.txt
	@echo "Result  Path                  Test Name                                 Phases" >> testreport.txt
	@echo "------  --------------------  ----------------------------------------  ------" >> testreport.txt
	for dir in $(DIRS); do make -C $$dir test;done
	@echo >> testreport.txt
	@echo "Report Complete: \c" >> testreport.txt
	@date >> testreport.txt
	@echo "*** TEST completed in test ***"
	@cat testreport.txt


clean:
	rm -f testreport.txt
	cd usersim;make clean

packageqa:
	cd usersim;make packageqa
	cd phonebookqa;make packageqa
