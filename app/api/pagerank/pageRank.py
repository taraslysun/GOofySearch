import numpy as np


def pagerank(adj_matrix, site_names, d=0.85, max_iter=100, tol=1e-6):
    # Number of pages
    n = adj_matrix.shape[0]
    
    # Create the personalization vector based on preferred categories
    personalization = np.ones(n) / n
    
    
    # Create transition matrix from adjacency matrix
    out_link_sums = adj_matrix.sum(axis=1)
    transition_matrix = adj_matrix / out_link_sums[:, np.newaxis]
    transition_matrix = np.nan_to_num(transition_matrix)  # Replace NaNs with 0 for pages with no out links

    # Initialize PageRank vector with the personalization vector
    pagerank = np.copy(personalization)
    
    # Iterative calculation of PageRank
    for _ in range(max_iter):
        new_pagerank = d * transition_matrix.T.dot(pagerank) + (1 - d) * personalization
        # Check for convergence
        if np.linalg.norm(new_pagerank - pagerank, 1) < tol:
            break
        pagerank = new_pagerank
    
    # Order sites by PageRank score
    ordered_sites = sorted(range(n), key=lambda i: -pagerank[i])
    ordered_sites = [(site_names[i], pagerank[i]) for i in ordered_sites]

    
    return ordered_sites

if __name__ == '__main__':
    # adj_matrix, site_categories, site_names = generate_site_data(10000)
    adj_matrix = np.loadtxt('adj_matrix.csv', dtype=int)
    site_names = np.loadtxt('site_names.csv', dtype=str)

    ordered_sites = pagerank(adj_matrix, site_names)
    print(ordered_sites)