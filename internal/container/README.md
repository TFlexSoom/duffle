# LinkedTree vs. GraphTree
We have all heard that Linked Lists are the bane of cacheing for languages. As that is the case
I created a basic implementation for a Tree (using linked references) and a Graph Implementation
(using a map for all relationships). This way, data can be sliced directly from the Graph Implementation
and I can test between the 2 implementations for speed. However, upon testing high depth 
data relationships, I found the tests ran pretty equivalently on my computer. I do have a top tier
machine, but I still expected more of a performance drop from a linked list which has to concatenate
all of the values when asked. Regardless, both implementations come with the same interface/contract
such that they can be used interchangeably with the IRs. Both models have also been provided with
custom capacities such that we can change them at runtime if desired.

## Notes And Desires
If there is further findings between the implementations, write them here