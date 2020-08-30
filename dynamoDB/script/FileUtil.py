import csv
import pprint

class FileUtil():

    #
    # 全データの取得
    # @return 値(配列)
    #
    def loadCSVData(self, file_path):
        with open(file_path) as f:
            reader = csv.DictReader(f)
            l = [row for row in reader]
        return l


