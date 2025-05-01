import matplotlib.pyplot as plt
import pandas as pd

file_path = 'benchmark-results.csv'

df = pd.read_csv(file_path)

plt.figure(figsize=(10, 6))
plt.bar(df['method'], df['ms_per_op'])
plt.xlabel('Method')
plt.ylabel('Milliseconds per Operation')
plt.title('Benchmark Results')
plt.xticks(rotation=45)
plt.tight_layout()
plt.savefig('benchmark_results.png')
