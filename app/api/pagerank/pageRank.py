import numpy as np


def pagerank(adj_matrix, site_categories, site_names, preferred_categories=[], d=0.85, max_iter=1000, tol=1e-6):
    n = adj_matrix.shape[0]
    
    personalization = np.zeros(n)
    for i in range(n):
        personalization[i] = 1
    
    if personalization.sum() > 0:
        personalization /= personalization.sum()
    else:
        personalization = np.ones(n) / n
    out_link_sums = adj_matrix.sum(axis=1)
    transition_matrix = adj_matrix / out_link_sums[:, np.newaxis]
    transition_matrix = np.nan_to_num(transition_matrix)

    pagerank = np.copy(personalization)
    
    for _ in range(max_iter):
        new_pagerank = d * transition_matrix.T.dot(pagerank) + (1 - d) * personalization
        if np.linalg.norm(new_pagerank - pagerank, 1) < tol:
            break
        pagerank = new_pagerank
    
    ordered_sites = sorted(range(n), key=lambda i: -pagerank[i])
    ordered_sites = [(site_names[i], pagerank[i], site_categories[i]) for i in ordered_sites]

    
    return ordered_sites

if __name__ == '__main__':
    preferred_categories = ['SCIENCE']
    adj_matrix = np.loadtxt('adj_matrix.csv', dtype=int)
    site_categories = np.loadtxt('site_categories.csv', dtype=str)
    site_names = np.loadtxt('site_names.csv', dtype=str)

    ordered_sites = pagerank(adj_matrix, site_categories, site_names, preferred_categories)

    for site, score in ordered_sites[:10]:
        print(f'{site}: {score}')
    print('\n')
