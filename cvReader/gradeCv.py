

# TODO: Implement the grading system
# 1. Read the cv
# 2. Tokenize the cv
# 3. Grade the cv based on the following criteria:
# - Number of pages
# - ATS Friendly
# - Check grammar
# - Check spelling
# - Core sections: Summary, Contact, Experience, Education, Skills, (Certifications, Projects, Languages, Hobbies
# - Action verbs
# - Education


def grade_cv(cv):

    # tokenize the cv
    cv = cv.lower()
    cv = cv.split()

    print(cv)