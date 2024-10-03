import os

from cvReader import read_pdf
from gradeCv import grade_cv


def main():
    # loop files in dir

    directory_path = 'cvDir'

    for filename in os.listdir(directory_path):
        if filename.endswith(".pdf"):
            file_path = os.path.join(directory_path, filename)
            text = read_pdf(file_path)
            print(text)
            print('-' * 50)
            grade_cv(text)


def return_cv_texts():
    # loop files in dir

    directory_path = 'cvDir'

    cv_texts = []

    for filename in os.listdir(directory_path):
        if filename.endswith(".pdf"):
            file_path = os.path.join(directory_path, filename)
            text = read_pdf(file_path)
            print(text)
            print('-' * 50)
            cv_texts.append(text)

    return cv_texts


if __name__ == '__main__':
    main()