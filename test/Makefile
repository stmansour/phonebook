# DIRS=usersim ws
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
	@./totals.sh
	@echo >> testreport.txt
	@echo "Report Complete: \c" >> testreport.txt
	@date >> testreport.txt
	@echo "*** TEST completed in test ***"
	@cat testreport.txt


clean:
	rm -f testreport.txt dbqq*.sql
	for dir in $(DIRS); do make -C $$dir clean;done

packageqa:
	cd usersim;make packageqa
	cd phonebookqa;make packageqa

package:
	echo "*** Package completed in test"secure:
	for dir in $(DIRS); do make -C $${dir} secure;done
	@rm -f config.json confdev.json confprod.json

secure:
	for dir in $(DIRS); do make -C $$dir secure;done
	@rm -rf config.json confdev.json confprod.json
