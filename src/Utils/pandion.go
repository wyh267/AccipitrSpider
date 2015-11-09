

package Utils



import (
	"fmt"
	"errors"
	"github.com/outmana/log4jzl"
)

const DB_PRE string = "./DB/"
const MAX_KEY_LEN int = 200
const RECORD_LEN int64 = 232

type NodeInfo struct{
	KeyOffset		int64	//KVDB的key的偏移量（每条数据*KeyOffset为真实偏移量）
	ValueOffset 	int64	//KVDB的Value的偏移量
	ValueLens		int64	//KVDB的Value的长度
	Value   		string	//KVDB的Value的真实值
}




type PandionKV struct{
	Name 	string
	IndexInfo map[string]NodeInfo
	indexMmap	*Mmap
	detailMmap	*Mmap
	MaxKeyOffset	int64
	Logger		*log4jzl.Log4jzl
}


/*****************************************************************************
*  function name : NewPadionKV
*  params : name string
*  return : *PandionKV
*
*  description : 通过名称新建一个KVDB
*
******************************************************************************/
func NewPandionKV(name string,logger *log4jzl.Log4jzl)(*PandionKV){
	
	
	
	this:=&PandionKV{MaxKeyOffset:0,Logger:logger,Name:name,indexMmap:nil,detailMmap:nil,IndexInfo:make(map[string]NodeInfo,0)}
	
	var err error
	idx_file_name:=fmt.Sprintf("%v/%v.idx",DB_PRE,this.Name)
	this.indexMmap,err=NewMmap(idx_file_name,MODE_CREATE)
	if err!=nil{
		return nil
	}
	this.indexMmap.AppendInt64(0)
	this.indexMmap.SetFileEnd(8)
	dtl_file_name:=fmt.Sprintf("%v/%v.dtl",DB_PRE,this.Name)
	this.detailMmap,err=NewMmap(dtl_file_name,MODE_CREATE)
	if err!=nil{
		return nil
	}
	this.detailMmap.SetFileEnd(0)
	return this

}




func NewPandionKVWithFile(name string,logger *log4jzl.Log4jzl)(*PandionKV){
	
	this:=&PandionKV{MaxKeyOffset:0,Logger:logger,Name:name,indexMmap:nil,detailMmap:nil,IndexInfo:make(map[string]NodeInfo,0)}
	
	this.loadLocalFiles()
	
	
	return this
}



func (this *PandionKV)loadLocalFiles()error{
	
	
	var err error
	idx_file_name:=fmt.Sprintf("%v/%v.idx",DB_PRE,this.Name)
	this.indexMmap, err = NewMmap(idx_file_name, MODE_APPEND)
	if err != nil {
		fmt.Printf("mmap error : %v \n", err)
		return err
	}
	
	
	dtl_file_name:=fmt.Sprintf("%v/%v.dtl",DB_PRE,this.Name)
	this.detailMmap, err = NewMmap(dtl_file_name, MODE_APPEND)
	if err != nil {
		fmt.Printf("mmap error : %v \n", err)
		return err
	}
	
	this.MaxKeyOffset=this.indexMmap.ReadInt64(0)
	
	var start int64 = 8
	var i int64=0
	//fmt.Printf("LEN: %v \n",idflen)
	for i=0;i<this.MaxKeyOffset;i++{
		
		offset:=i*RECORD_LEN+8
		
		keyoffset := this.indexMmap.ReadInt64(offset+200)
		valueoffset := this.indexMmap.ReadInt64(offset+208)
		valuelen := this.indexMmap.ReadInt64(offset+216)
		keylen := this.indexMmap.ReadInt64(offset+224)
		key := this.indexMmap.ReadString(offset,keylen)
		
		this.IndexInfo[key]=NodeInfo{KeyOffset:keyoffset,ValueOffset:valueoffset,ValueLens:valuelen,Value:""}
		start+=RECORD_LEN
	}
	this.indexMmap.SetFileEnd(start)
	
	
	return nil
	
	
	
	
	
}



func (this *PandionKV)GetAllKeys()(map[string][]string,error){
	

	rs:=make(map[string][]string)
	for k,v := range this.IndexInfo{
		_,ok:=rs[v.Value]
		if !ok {
			rs[v.Value]=make([]string,0)
			rs[v.Value]=append(rs[v.Value],k)
		}else{
			rs[v.Value]=append(rs[v.Value],k)
		}
		
	}
	
	return rs,nil
}


func (this *PandionKV)getFromDisk(offset,lens int64) (string,error){
	
	return this.detailMmap.ReadString(offset,lens),nil
	
}



func (this *PandionKV)Set(key,value string) error{
	
	if len(key) > MAX_KEY_LEN{
		//this.Logger.Printf("[ERROR] Key is too long")
		return errors.New("Key is too long")
	}
	var maxkeyoffset int64 = this.MaxKeyOffset*RECORD_LEN + 8
	v,ok:=this.IndexInfo[key]
	if !ok {
		//新建一个节点
		this.IndexInfo[key]=NodeInfo{KeyOffset:maxkeyoffset,ValueOffset:this.detailMmap.GetPointer(),ValueLens:int64(len(value)),Value:value}
		this.MaxKeyOffset++
		this.addToDisk(key,true)
		return nil
		
	}else{
		//修改一个节点
		v.Value=value
		if v.ValueLens != int64(len(value)){
			v.ValueOffset=this.detailMmap.GetPointer()
			v.ValueLens=int64(len(value))
			this.IndexInfo[key]=v
			this.addToDisk(key,true)
		}else{
			this.IndexInfo[key]=v
			this.addToDisk(key,false)
		}
		
	}
	
	return nil
}




func (this *PandionKV)Get(key string) (string,error){
	
	
	v,ok:=this.IndexInfo[key]
	if !ok {
		return "",errors.New("Key is not found")
	}
	
	if v.Value ==""{	
		value,err :=this.getFromDisk(v.ValueOffset,v.ValueLens)
		if err!=nil{
			return "",err
		}
		v.Value=value
		this.IndexInfo[key]=v
	}
	
	return v.Value,nil
}


/*
idx ..
| lens(8bytes) |
| key (200bytes) | KeyOffset(8bytes) | ValueOffset(8bytes) | ValueLens (8bytes)| KeyLen(8bytes) |
| key (200bytes) | KeyOffset(8bytes) | ValueOffset(8bytes) | ValueLens (8bytes)| KeyLen(8bytes) |
| key (200bytes) | KeyOffset(8bytes) | ValueOffset(8bytes) | ValueLens (8bytes)| KeyLen(8bytes) |
| key (200bytes) | KeyOffset(8bytes) | ValueOffset(8bytes) | ValueLens (8bytes)| KeyLen(8bytes) |

dtl ..
| value1 (...) | value2 (...) | value3 (...) | value4 (...) |

*/

func (this *PandionKV)addToDisk(key string,isappend bool) error {
	
	v,_:=this.IndexInfo[key]
	this.indexMmap.WriteInt64(0,this.MaxKeyOffset)
	fmt.Printf("Max : %v \n",this.MaxKeyOffset)
	start:=v.KeyOffset
	
	str_bytes := make([]byte, MAX_KEY_LEN)
	copy(str_bytes, []byte(key))
	this.indexMmap.WriteBytes(start,str_bytes)
	this.indexMmap.WriteInt64(start+int64(MAX_KEY_LEN),v.KeyOffset)
	this.indexMmap.WriteInt64(start+int64(MAX_KEY_LEN)+8,v.ValueOffset)
	this.indexMmap.WriteInt64(start+int64(MAX_KEY_LEN)+16,v.ValueLens)
	this.indexMmap.WriteInt64(start+int64(MAX_KEY_LEN)+24,int64(len(key)))
	//this.detailMmap.WriteString(v.ValueOffset,v.Value)
	if isappend{
		this.detailMmap.AppendString(v.Value)
	}else{
		this.detailMmap.WriteString(v.ValueOffset,v.Value)
	}
	
	return nil
	
}











