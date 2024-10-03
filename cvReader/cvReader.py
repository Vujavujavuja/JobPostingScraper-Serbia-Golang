import pypdf
import re
import nltk
from nltk.corpus import words

english_words = set(words.words())


def read_pdf(file_path):
    print(file_path)
    pdf = pypdf.PdfReader(file_path)
    text = ''
    for page in pdf.pages:
        text += page.extract_text()
    print(text)
    text = fix_spaces(text)
    text = fix_spaces_with_dictionary(text)
    return text


def fix_spaces(text):
    processed_text = re.sub(r'([a-z])([A-Z])', r'\1 \2', text)
    processed_text = re.sub(r'([,.!?])([A-Za-z])', r'\1 \2', processed_text)
    processed_text = re.sub(r'([;:])([A-Za-z])', r'\1 \2', processed_text)

    return processed_text


def fix_spaces_with_dictionary(text):
    words_in_text = text.split()
    corrected_text = []

    for word in words_in_text:
        if word.lower() in english_words:
            corrected_text.append(word)
        else:
            corrected_text.append(split_word_using_dictionary(word))

    return ' '.join(corrected_text)


def split_word_using_dictionary(word):
    for i in range(1, len(word)):
        left, right = word[:i], word[i:]
        if left.lower() in english_words and right.lower() in english_words:
            return left + ' ' + right

    return word