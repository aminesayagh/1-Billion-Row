import matplotlib.pyplot as plt
import networkx as nx

G = nx.DiGraph()

# Adding nodes with labels
G.add_node("q0", label="Start")
G.add_node("q1", label="-/+")
G.add_node("q2", label="1..9")
G.add_node("q3", label=".")
G.add_node("q4", label="0..9")
G.add_node("q5", label="1..9")
G.add_node("q6", label="Accept")    



# Adding edges with labels
G.add_edge("q0", "q1", label="+/-")
G.add_edge("q0", "q2", label="1..9")
G.add_edge("q1", "q2", label="1..9")
G.add_edge("q2", "q2", label="0..9")
G.add_edge("q2", "q3", label=".")
G.add_edge("q3", "q4", label="0")
G.add_edge("q3", "q5", label="1..9")
G.add_edge("q4", "q4", label="0..9")
G.add_edge("q4", "q5", label="1..9")
G.add_edge("q5", "q4", label="0")
G.add_edge("q5", "q5", label="1..9")
G.add_edge("q5", "q6", label="Accept")

# Node positions
pos = {
    "q0": (0, 1),
    "q1": (2, 2),
    "q2": (2, 0),
    "q3": (4, 0),
    "q4": (6, 1),
    "q5": (6, -1),
    "q6": (4, -1)
}


# Draw the nodes and edges
nx.draw_networkx_nodes(G, pos, node_size=700, node_color="lightblue")
nx.draw_networkx_edges(G, pos, arrowsize= 16)
nx.draw_networkx_labels(G, pos, font_size=9, font_weight="bold")

# Draw edge labels
edge_labels = nx.get_edge_attributes(G, 'label')
nx.draw_networkx_edge_labels(G, pos, edge_labels=edge_labels, font_size=9, rotate=False)

# Save the graph as SVG and PNG
plt.savefig("automaton.svg")
plt.savefig("automaton.png")

# Show the diagram
plt.show()
