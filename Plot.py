import yaml
import numpy as np
import matplotlib.pyplot as plt

with open('./result.yaml') as f:
    data = yaml.load_all(f, Loader=yaml.FullLoader)
    for d in data:
        fileNum = d['filenum']
        rawEntropy = d['entropy']['rawentropy']
        estimatedEntropy = d['entropy']['estimatedentropy']
        compressedEntropy = d['entropy']['compressedentropy']

x = np.arange(len(fileNum))
width = 0.25

fig, ax = plt.subplots()
rect1 = ax.bar(x-width, rawEntropy, width, label='Raw Entropy')
rect2 = ax.bar(x, estimatedEntropy, width, label='Estimated Entropy')
rect3 = ax.bar(x+width, compressedEntropy, width, label='Compressed Entropy')

ax.set_xticks(x)
ax.set_xticklabels(fileNum)
ax.legend()
fig.tight_layout()

plt.show()