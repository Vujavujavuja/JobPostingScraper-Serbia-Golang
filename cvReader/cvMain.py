from cvReader import read_pdf


def main():
    file_path = 'cvDir/NemanjaVujicCV.pdf'
    text = read_pdf(file_path)
    print(text)


if __name__ == '__main__':
    main()