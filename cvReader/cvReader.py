import pytesseract
from pdf2image import convert_from_path
import re
import nltk
from nltk.corpus import words

# Download the word corpus if not already available
nltk.download('words')

# Set path to tesseract executable
pytesseract.pytesseract.tesseract_cmd = r'C:\Program Files\Tesseract-OCR\tesseract.exe'

# Set path to Poppler binaries
poppler_path = r'C:\Program Files\poppler-24.07.0\Library\bin'

# Load the English word list
english_words = set(words.words())


def read_pdf(file_path):
    # Convert PDF pages to images using pdf2image and specifying the Poppler path
    images = convert_from_path(file_path, poppler_path=poppler_path)

    text = ''
    for image in images:
        # Extract text from each image using Tesseract OCR
        text += pytesseract.image_to_string(image)

    # Optional: Process the text after extraction
    text = fix_spaces(text)
    text = fix_spaces_with_dictionary(text)

    return text


def fix_spaces(text):
    """Fix common spacing issues between words and punctuation."""
    processed_text = re.sub(r'([a-z])([A-Z])', r'\1 \2', text)
    processed_text = re.sub(r'([,.!?])([A-Za-z])', r'\1 \2', processed_text)
    processed_text = re.sub(r'([;:])([A-Za-z])', r'\1 \2', processed_text)
    return processed_text


def fix_spaces_with_dictionary(text):
    """Try to fix spacing issues using a dictionary-based approach."""
    words_in_text = text.split()
    corrected_text = []

    for word in words_in_text:
        if word.lower() in english_words:
            corrected_text.append(word)
        else:
            corrected_text.append(split_word_using_dictionary(word))

    return ' '.join(corrected_text)


def split_word_using_dictionary(word):
    """Attempt to split a word into valid English words."""
    for i in range(1, len(word)):
        left, right = word[:i], word[i:]
        if left.lower() in english_words and right.lower() in english_words:
            return left + ' ' + right
    return word
