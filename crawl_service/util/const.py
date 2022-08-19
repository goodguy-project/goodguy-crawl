import os
from os.path import dirname

from crawl_service import crawl_service_pb2

ROOT = dirname(dirname(dirname(os.path.abspath(__file__))))

PLATFORM_RECENT_CONTEST = [
    'atcoder',
    'codeforces',
    'leetcode',
    'luogu',
    'nowcoder',
    'codechef',
    'acwing',
]

PLATFORM_ALL = [
    'atcoder',
    'codeforces',
    'leetcode',
    'luogu',
    'nowcoder',
    'vjudge',
    'codechef',
    'acwing',
]

VJUDGE_VERDICT_MAP = {
    'AC': crawl_service_pb2.Accepted,
    'WA': crawl_service_pb2.WrongAnswer,
    'CE': crawl_service_pb2.CompilationError,
    'TLE': crawl_service_pb2.TimeLimitExceeded,
    'MLE': crawl_service_pb2.MemoryLimitExceeded,
    'RE': crawl_service_pb2.RuntimeError,
    'OLE': crawl_service_pb2.OutputLimitExceeded,
    'PE': crawl_service_pb2.PresentationError,
}

VJUDGE_PROGRAMMING_LANGUAGE_MAP = {
    'CPP': crawl_service_pb2.Cpp,
    'C': crawl_service_pb2.C,
    'JAVA': crawl_service_pb2.Java,
    'PYTHON': crawl_service_pb2.Python,
    'RUBY': crawl_service_pb2.Ruby,
    'CSHARP': crawl_service_pb2.CSharp,
}

CODEFORCES_VERDICT_MAP = {
    'OK': crawl_service_pb2.Accepted,
    'COMPILATION_ERROR': crawl_service_pb2.CompilationError,
    'RUNTIME_ERROR': crawl_service_pb2.RuntimeError,
    'WRONG_ANSWER': crawl_service_pb2.WrongAnswer,
    'PRESENTATION_ERROR': crawl_service_pb2.PresentationError,
    'TIME_LIMIT_EXCEEDED': crawl_service_pb2.TimeLimitExceeded,
    'MEMORY_LIMIT_EXCEEDED': crawl_service_pb2.MemoryLimitExceeded,
    'IDLENESS_LIMIT_EXCEEDED': crawl_service_pb2.IdlenessLimitExceeded,
    'SECURITY_VIOLATED': crawl_service_pb2.SecurityViolated,
    'CRASHED': crawl_service_pb2.Crashed,
    'INPUT_PREPARATION_CRASHED': crawl_service_pb2.InputPreparationCrashed,
    'CHALLENGED': crawl_service_pb2.Challenged,
    'SKIPPED': crawl_service_pb2.Skipped,
    'TESTING': crawl_service_pb2.Testing,
    'REJECTED': crawl_service_pb2.Rejected,
}
