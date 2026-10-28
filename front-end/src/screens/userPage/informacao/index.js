import React, { useState, useEffect } from "react";
import { View, Text, TouchableOpacity, TextInput, Image, ActivityIndicator, FlatList } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import { useNavigation } from "@react-navigation/native";

import style from "./style";
import { listarAnimais } from "../../../../services/API"; 

export default () => {
  const [animais, setAnimais] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [buscaTexto, setBuscaTexto] = useState("");
  
  const navigation = useNavigation();

  const carregarAnimais = async () => {
    try {
      setLoading(true);
      setError(null);
      const dados = await listarAnimais();
      setAnimais(dados);
      console.log(dados)
    } catch (err) {
      console.error("Erro ao carregar os animais:", err);
      setError("Não foi possível carregar os animais.");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    carregarAnimais();
  }, []);

  const animaisFiltrados = animais.filter((item) => {
    const termo = buscaTexto.toLowerCase();
    const nome = (item.nome || item.nome_popular || "").toLowerCase();
    const descricao = (item.descricao || item.especie || "").toLowerCase();

    return nome.includes(termo) || descricao.includes(termo);
  });

    const renderCardItem = ({ item }) => (
    <TouchableOpacity
      style={style.cardContainer}
      onPress={() => navigation.navigate("animalDetalhe", { id: item.id, resumo: item })}
      activeOpacity={0.8}
    >
      <Image
        source={{ uri: item.imagem_url || "https://via.placeholder.com/300x180?text=Sem+Imagem" }}
        style={style.cardImage}
        resizeMode="cover"
      />
      <View style={style.cardTextArea}>
        <Text style={style.cardTitle}>{item.nome || item.nome_popular || "Sem nome"}</Text>
        <Text style={style.cardDescription} numberOfLines={3}>
          {item.breve_descricao || item.descricao || item.especie || "Sem descrição disponível."}
        </Text>
      </View>
    </TouchableOpacity>
  );


  return (
    <SafeAreaView style={style.SafeAreaView}>
      <View style={style.busca}>
        <TextInput 
          style={style.buscaInput} 
          placeholder="Busca por nome ou espécie..." 
          value={buscaTexto}                
          onChangeText={setBuscaTexto}      
          autoCapitalize="none"
        />
      </View>

      {/* Conteúdo Dinâmico */}
      {loading ? (
        <View style={style.centerContainer}>
          <ActivityIndicator size="large" color="#007AFF" />
          <Text style={style.loadingText}>Carregando animais...</Text>
        </View>
      ) : error ? (
        <View style={style.centerContainer}>
          <Text style={style.errorText}>{error}</Text>
          <TouchableOpacity style={style.retryButton} onPress={carregarAnimais}>
            <Text style={style.retryButtonText}>Tentar Novamente</Text>
          </TouchableOpacity>
        </View>
      ) : (
        <View style={{ flex: 1, width: "100%", marginTop: 10 }}>
          <FlatList
            data={animaisFiltrados}
            keyExtractor={(item, index) => item.id ? item.id.toString() : index.toString()}
            renderItem={renderCardItem}
            showsVerticalScrollIndicator={false}
            ListEmptyComponent={
              <View style={style.centerContainer}>
                <Text style={style.emptyText}>Nenhum animal encontrado.</Text>
              </View>
            }
          />
        </View>
      )}
    </SafeAreaView>
  );
};