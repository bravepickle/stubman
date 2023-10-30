{{ define "title"}}<title>Stubman | Stubs List</title>{{ end }}
{{ define "content"}}<h2>List</h2>
<a href="{{ .BaseUri }}/stubman/create/" class="btn btn-success">Create</a>
<table class="table table-striped">
<thead>
	<tr>
		<th>ID</th>
		<th>Name</th>
		<th>Method</th>
		<th>URI</th>
		<th>Created</th>
		<th>Last Viewed</th>
		<th>Views</th>
		<th>Actions</th>
	</tr>
</thead>
<tbody>
{{range .Data}}
    <tr>
		<td>{{.Id}}</td>
		<td>{{.Name}}</td>
		<td>{{.RequestMethod}}</td>
		<td>{{.RequestUri}}</td>
		<td>{{.Created}}</td>
		<td>{{.LastViewed}}</td>
		<td>{{.Views}}</td>
		<td>
			<a href="edit/{{.Id}}" title="edit"><span class="glyphicon glyphicon-pencil" aria-hidden="true"></span></a>
			<a href="#" data-href="{{ $.BaseUri }}/stubman/delete/{{.Id}}" data-toggle="modal" data-target="#del-confirm" class="btn-del" title="delete"><span class="glyphicon glyphicon-remove" aria-hidden="true"></span></a>
			<a href="copy/{{.Id}}" title="copy"><span class="glyphicon glyphicon-copy" aria-hidden="true"></span></a>
		</td>
	</tr>
{{end}}
</tbody>
</table>

<!-- Modal -->
<div class="modal fade" id="del-confirm" tabindex="-1" role="dialog" aria-labelledby="myModalLabel">
  <div class="modal-dialog" role="document">
    <div class="modal-content">
      <div class="modal-header">
        <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
        <h4 class="modal-title" id="myModalLabel">Deletion Confirmation</h4>
      </div>
      <div class="modal-body">
        Are you sure you want to delete this record?
      </div>
      <div class="modal-footer">
        <button type="button" class="btn btn-default" id="btn-del-cancel" data-dismiss="modal">Cancel</button>
        <button type="button" class="btn btn-danger" id="btn-del-confirm">Delete</button>
      </div>
    </div>
  </div>
</div>

{{ end }}

{{ define "scripts" }}
<script src="{{ .BaseUri }}/stubman/static/js/index.js" crossorigin="anonymous"></script>
<script type="text/javascript">
   $(document).ready(function() {
		initIndex();
   });
</script>
{{ end }}
