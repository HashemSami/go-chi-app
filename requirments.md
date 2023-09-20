1. create a gallery with a title
2. upload images to a gallery
3. delete images from a gallery
4. update the title of a gallery
5. view a gallery (so we can share it with others).
6. delete a gallery
7. view a list of galleries we are allowed to edit

todo that we need:

views:

- create a new gallery
- edit a gallery
- view an existing gallery
- view a list of all of our galleries

new, edit, show, index

we will also need controllers (aka HTTP handlers) to support these views:

handlers:

- New and Create to render and process a new gallery form
- Edit and Update to render and process a new gallery form
- Show to render a gallery
- Delete to delete a gallery

- an HTTP handler to process image uploads
- an HTTP handler to remove images from a gallery

finally, we need a way to persist data in our models package, and this will need to support the following:

- create a gallery
- update a gallery
- querying a gallery
- querying for all galleries with a user ID
- deleting a gallery
- creating an image for a gallery
- deleting an image from a gallery
