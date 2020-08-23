#!/usr/local/bin/python3

import re
import string
import os
import sys

import yake

# Input config
cwd = sys.argv[1]
file = sys.argv[2]

file_name, file_ext = os.path.splitext(os.path.basename(file))

title = file_name
destination = os.path.join(cwd, "keywords")
meta_destination = os.path.join(destination, ".meta")
meta_file = os.path.join(meta_destination, title + ".kwds")

# yake config
language = "en"
max_ngram_size = 3
deduplication_thresold = 0.75
deduplication_algo = "jaro"
window_size = 1

# Templates
strip_punct_re = r"[{}]".format(string.punctuation)

start_auto_tags = "<!-- start auto-tags -->"
end_auto_tags = "<!-- end auto-tags -->"

backlink_template = string.Template("[[$to]]\n")
file_template = string.Template(
    "# Auto-tag: $tag\n\n" + start_auto_tags + "\n"
)
tag_file_path_template = string.Template(destination + "/$tag.md")

# Create keyword folder
if not os.path.exists(destination):
    os.makedirs(destination)

# Create meta folder
if not os.path.exists(meta_destination):
    os.makedirs(meta_destination)


# Helpers
def load_prev_keywords():
    if not os.path.exists(meta_file):
        return []

    with open(meta_file, "r+", encoding="utf-8") as file_object:
        return [line.strip() for line in file_object.readlines()]


def save_keywords(keywords):
    with open(meta_file, "w+", encoding="utf-8") as file_object:
        file_object.write("\n".join(keywords))


def diff_keyword_lists(new):
    old_set = set(load_prev_keywords())
    new_set = set(new)

    removed = old_set - new_set
    added = new_set - old_set

    return list(removed), list(added)


def add_backlink(tag):
    path = tag_file_path_template.substitute(tag=tag)

    if not os.path.exists(path):
        with open(path, "w+", encoding="utf-8") as file_object:
            file_object.write(
                file_template.substitute(tag=tag) +
                backlink_template.substitute(to=title) +
                end_auto_tags
            )
    else:
        with open(path, "r") as file_object:
            lines = file_object.readlines()

        # Add the new backlink above the end_auto_tags marker
        with open(path, "w+", encoding="utf-8") as file_object:
            for line in lines:
                if line == end_auto_tags:
                    file_object.write(backlink_template.substitute(to=title))

                file_object.write(line)


def remove_backlink(tag):
    path = tag_file_path_template.substitute(tag=tag)
    file_has_links = False

    with open(path, "r") as file_object:
        lines = file_object.readlines()

    # Filter out backlinks to the current file
    with open(path, "w") as file_object:
        for line in lines:
            if line != backlink_template.substitute(to=title):
                file_object.write(line)

                if line.startswith("[["):
                    file_has_links = True

    # If we didn't encounter any backlinks besides the one we filtered,
    # delete the now empty file
    if not file_has_links:
        os.remove(path)


def process_file():
    with open(file, "r+", encoding="utf-8") as content_file:
        contents = re.sub(
            r"\[\/\/begin[\s\S]+\/\/end\][\s\S]*$", "", content_file.read())
        length = len(contents.split())

        kw_extractor = yake.KeywordExtractor(
            lan=language,
            n=max_ngram_size,
            dedupLim=deduplication_thresold,
            dedupFunc=deduplication_algo,
            windowsSize=window_size,
            top=max(3, min(round(length * 0.03), 10)),
            features=None,
        )

        keywords = kw_extractor.extract_keywords(
            re.sub(r"[’‘]", "'", contents)
        )

        tags = [re.sub(strip_punct_re, "", w).replace(" ", "-")
                for w, score in keywords]

        removed, added = diff_keyword_lists(tags)

        for to_remove in removed:
            print("Removing link from " + to_remove + " to " + title)
            remove_backlink(to_remove)

        for to_add in added:
            print("Adding link from " + to_add + " to " + title)
            add_backlink(to_add)

        save_keywords(tags)
        print("Processing complete")


if file_ext == ".md":
    print("Processing " + title + file_ext)
    process_file()
