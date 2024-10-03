import pypdf


def read_pdf(file_path):
    print(file_path)
    pdf = pypdf.PdfReader(file_path)
    text = ''
    for page in pdf.pages:
        text += page.extract_text()
    print(text)
    return text
