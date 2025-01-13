import sys
import os
import json
from git import Repo
import re

def analyze_file(file_path, skills):
    with open(file_path, 'r') as file:
        lines = file.readlines()
        
    results = {skill['id']: {'score': 0, 'evidence': []} for skill in skills}
    
    for i, line in enumerate(lines, 1):
        for skill in skills:
            if skill['id'] == 'junit_tests' and '@Test' in line:
                results[skill['id']]['score'] = 1
                results[skill['id']]['evidence'].append(i)
            elif skill['id'] == 'mockito_tests' and 'import org.mockito' in line:
                results[skill['id']]['score'] = 1
                results[skill['id']]['evidence'].append(i)
    
    return results

def get_file_author(repo, file_path):
    blame = repo.git.blame('-p', file_path)
    author_line = re.search(r'author (.*)', blame)
    if author_line:
        return author_line.group(1)
    return "Unknown"

def get_commit_id(repo, file_path):
    return repo.git.log('-n', '1', '--pretty=format:%H', '--', file_path)

def analyze_repo(repo_path, files_to_analyze):
    repo = Repo(repo_path)
    
    skills = [
        {"id": "junit_tests", "name": "Writing unit tests using JUnit"},
        {"id": "mockito_tests", "name": "Writing unit tests using Mockito"}
    ]
    
    if not files_to_analyze:
        files_to_analyze = [
            item.path for item in repo.tree().traverse()
            if item.path.endswith('.java')
        ]
    
    results = {}
    
    for file_path in files_to_analyze:
        full_path = os.path.join(repo_path, file_path)
        if os.path.exists(full_path) and file_path.endswith('.java'):
            author = get_file_author(repo, file_path)
            commit_id = get_commit_id(repo, file_path)
            file_results = analyze_file(full_path, skills)
            
            if author not in results:
                results[author] = {skill['id']: {'score': 0, 'evidence': []} for skill in skills}
            
            for skill_id, data in file_results.items():
                if data['score'] > results[author][skill_id]['score']:
                    results[author][skill_id]['score'] = data['score']
                    results[author][skill_id]['evidence'] = []  # Reset evidence if we found a higher score
                
                if data['score'] == results[author][skill_id]['score']:
                    results[author][skill_id]['evidence'].extend([
                        {"commit_id": commit_id, "file": file_path, "line": line}
                        for line in data['evidence']
                    ])
    
    output = []
    for author, scores in results.items():
        for skill in skills:
            skill_id = skill['id']
            output.append({
                "metric_id": skill_id,
                "git_author": author,
                "score": str(scores[skill_id]['score']),
                "evidence": scores[skill_id]['evidence']
            })
    
    return json.dumps(output, indent=2)

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: python script.py <repo_path> [file1] [file2] ...")
        sys.exit(1)
    
    repo_path = sys.argv[1]
    files_to_analyze = sys.argv[2:] if len(sys.argv) > 2 else []
    
    result = analyze_repo(repo_path, files_to_analyze)
    print(result)