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

### Discovery From 12/19/2023
One thing I do realize is with a single base object (the graph/map implementation) you can capture the capacity
for your tree in a singular struct which would be copied either by reference or otherwise in a LinkedListMode.
__Maybe we need to do some reframing solution so it can be done with linkedlists__ At least in theory we are
getting a memory savings of some sort with the Graph Implementation.