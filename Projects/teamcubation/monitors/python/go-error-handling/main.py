import os
import json
import subprocess
import re
from collections import defaultdict
from typing import List, Dict, Any

from git import Repo, InvalidGitRepositoryError

def get_git_author(commit_hash: str, repo_path: str) -> str:
    try:
        result = subprocess.run(
            ['git', 'show', '-s', '--format=%an <%ae>', commit_hash],
            cwd=repo_path,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True
        )
        return result.stdout.strip()
    except Exception as e:
        print(f"Error getting git author: {e}")
        return "Unknown"

def analyze_panic_usage(file_content: str, file_path: str) -> List[Dict[str, Any]]:
    evidence = []
    if "cmd" in file_path:
        return evidence

    for i, line in enumerate(file_content.splitlines(), start=1):
        if "panic(" in line:
            evidence.append({
                "file": file_path,
                "line": i
            })
    return evidence

def analyze_error_wrapping(file_content: str, file_path: str) -> List[Dict[str, Any]]:
    evidence = []
    error_wrap_patterns = [r"fmt\.Errorf\(", r"errors\.Wrap\("]

    for i, line in enumerate(file_content.splitlines(), start=1):
        for pattern in error_wrap_patterns:
            if re.search(pattern, line):
                evidence.append({
                    "file": file_path,
                    "line": i
                })
    return evidence

def analyze_errors_ignore(file_content: str, file_path: str) -> List[Dict[str, Any]]:
    evidence = []
    lines = file_content.splitlines()

    for i, line in enumerate(lines, start=1):
        if re.search(r'_,\s*err\s*:=', line):
            if i == len(lines) or not re.search(r'if\s+err\s*!=\s*nil', lines[i]):
                evidence.append({
                    "file": file_path,
                    "line": i
                })
    return evidence

def analyze_repo(repo_path: str, modified_files: List[str]) -> List[Dict[str, Any]]:
    if not modified_files:
        modified_files = []
        for root, _, files in os.walk(repo_path):
            for file in files:
                if file.endswith(".go"):
                    modified_files.append(os.path.join(root, file))

    results = defaultdict(list)

    for file_path in modified_files:
        with open(file_path, 'r') as file:
            content = file.read()

        panic_evidence = analyze_panic_usage(content, file_path)
        wrap_evidence = analyze_error_wrapping(content, file_path)
        ignore_evidence = analyze_errors_ignore(content, file_path)

        if panic_evidence or wrap_evidence or ignore_evidence:
            # Solo obtener commit_hash y author si se detecta algún problema
            commit_hash = get_last_commit_hash(file_path, repo_path)
            author = get_git_author(commit_hash, repo_path)

            if panic_evidence:
                results["panic_usage"].append({
                    "git_author": author,
                    "score": "2",
                    "evidence": [{"commit_id": commit_hash, "file": e["file"], "line": e["line"]} for e in panic_evidence]
                })

            if not wrap_evidence:
                results["error_wrap"].append({
                    "git_author": author,
                    "score": "2",
                    "evidence": [{"commit_id": commit_hash, "file": e["file"], "line": e["line"]} for e in wrap_evidence]
                })

            if ignore_evidence:
                results["errors_ignore"].append({
                    "git_author": author,
                    "score": "1",
                    "evidence": [{"commit_id": commit_hash, "file": e["file"], "line": e["line"]} for e in ignore_evidence]
                })

    return results

def get_relative_path(repo_path, file_path):
    return os.path.relpath(file_path, repo_path)

def get_last_commit_hash(file_path: str, repo_path: str) -> str:
    try:
        repo = Repo(repo_path)

        # Verifica si el archivo está rastreado por el repositorio
        if get_relative_path(repo_path, file_path) not in repo.git.ls_files():
            raise FileNotFoundError(f"The file '{file_path}' is not tracked by Git in the repository.")

        # Obtiene el commit más reciente que modificó el archivo
        commits = list(repo.iter_commits(paths=file_path, max_count=1))
        if commits:
            return commits[0].hexsha
        else:
            return None
    except InvalidGitRepositoryError:
        print(f"The path '{repo_path}' is not a valid Git repository.")
        return None

def main(repo_path: str, modified_files: List[str]):
    # repo_path = '/home/osalomon/Projects/meli/tech-booster/go-basics'
    results = analyze_repo(repo_path, modified_files)
    print(json.dumps(results, indent=4))

if __name__ == "__main__":
    import sys
    repo_path = sys.argv[1]
    modified_files = sys.argv[2:]

    main(repo_path, modified_files)
