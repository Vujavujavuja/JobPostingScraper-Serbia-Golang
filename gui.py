import sqlite3
from tkinter import *
from tkinter import ttk
from tkinter import messagebox


# Connect to SQLite database (create if it doesn't exist)
def connect_db():
    conn = sqlite3.connect('jobs.db')  # Use your actual database name here
    cursor = conn.cursor()
    cursor.execute('''CREATE TABLE IF NOT EXISTS jobs (
                        id INTEGER PRIMARY KEY AUTOINCREMENT,
                        title TEXT,
                        company TEXT,
                        location TEXT,
                        seniority TEXT,
                        url TEXT UNIQUE,
                        Site TEXT
                    )''')
    conn.commit()
    conn.close()


# Fetch data based on search criteria
def search_data():
    conn = sqlite3.connect('jobs.db')
    cursor = conn.cursor()

    seniority = seniority_var.get()
    query = search_entry.get()

    if seniority == "All":
        cursor.execute(f"SELECT * FROM jobs WHERE title LIKE ? OR company LIKE ?",
                       ('%' + query + '%', '%' + query + '%'))
    else:
        cursor.execute(f"SELECT * FROM jobs WHERE (title LIKE ? OR company LIKE ?) AND seniority = ?",
                       ('%' + query + '%', '%' + query + '%', seniority))

    rows = cursor.fetchall()
    update_treeview(rows)
    conn.close()


# Clear search results and reload all data
def clear_data():
    search_entry.delete(0, END)
    load_data()


# Update Treeview with data
def update_treeview(rows):
    for i in tree.get_children():
        tree.delete(i)
    for row in rows:
        tree.insert("", END, values=row)


# Load all data from the SQLite database
def load_data():
    conn = sqlite3.connect('jobs.db')
    cursor = conn.cursor()
    cursor.execute("SELECT * FROM jobs")
    rows = cursor.fetchall()
    update_treeview(rows)
    conn.close()


# GUI setup
root = Tk()
root.title("Job Listings Viewer")
root.geometry("800x400")

# Seniority dropdown menu (preset options)
seniority_var = StringVar()
seniority_label = Label(root, text="Seniority:")
seniority_label.grid(row=0, column=0, padx=10, pady=10)

# Preset dropdown with Senior, Mid, Junior, Intern options
seniority_dropdown = ttk.Combobox(root, textvariable=seniority_var, values=["All", "Senior", "Mid", "Junior", "Intern"],
                                  state="readonly")
seniority_dropdown.grid(row=0, column=1, padx=10, pady=10)
seniority_dropdown.current(0)

# Search bar
search_label = Label(root, text="Search:")
search_label.grid(row=0, column=2, padx=10, pady=10)

search_entry = Entry(root)
search_entry.grid(row=0, column=3, padx=10, pady=10)

# Search and Clear buttons
search_button = Button(root, text="Search", command=search_data)
search_button.grid(row=0, column=4, padx=10, pady=10)

clear_button = Button(root, text="Clear", command=clear_data)
clear_button.grid(row=0, column=5, padx=10, pady=10)

# Treeview to display data
tree = ttk.Treeview(root, columns=("ID", "Title", "Company", "Location", "Seniority", "URL", "Site"), show="headings")
tree.heading("ID", text="ID")
tree.heading("Title", text="Title")
tree.heading("Company", text="Company")
tree.heading("Location", text="Location")
tree.heading("Seniority", text="Seniority")
tree.heading("URL", text="URL")
tree.heading("Site", text="Site")
tree.grid(row=1, column=0, columnspan=6, padx=10, pady=10)

# Scrollbar for Treeview
scrollbar = ttk.Scrollbar(root, orient="vertical", command=tree.yview)
scrollbar.grid(row=1, column=6, sticky="ns")
tree.configure(yscrollcommand=scrollbar.set)

# Load initial data from the database
connect_db()
load_data()

root.mainloop()
